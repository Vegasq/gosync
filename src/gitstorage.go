package main

import (
	"fmt"
	ssh2 "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

type GitStorage struct {
	repo *git.Repository
}

func (gs *GitStorage) openRepo() {
	r, err := git.PlainOpen(storageDir)
	if err != nil {
		fmt.Printf("Failed to open repo %s", err)
	}

	gs.repo = r
}

func getSshKeyAuth() transport.AuthMethod {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	privateSshKeyFile := usr.HomeDir + "/.ssh/id_rsa"

	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := ssh2.ParsePrivateKey([]byte(sshKey))
	auth = &ssh.PublicKeys{User: "git", Signer: signer}

	return auth
}

func repoIsEmptyError(err error) bool {
	if err.Error() == "remote repository is empty" {
		return true
	}
	return false
}

func repoAlreadyExistsError(err error) bool {
	if err.Error() == "repository already exists" {
		return true
	}
	return false
}

func createSignature() *object.Signature {
	return &object.Signature{
		Name:  "gosync",
		Email: "gosync",
		When:  time.Time{},
	}
}

func createAndPushEmptyRepo(){
	// Create repo
	repo, err := git.PlainInit(storageDir, false)
	if err != nil {
		fmt.Println("Failed to init repo")
		os.Exit(1)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{CONFIG.Main.Repo},
	})

	if err != nil {
		fmt.Println("Failed to create remote")
		os.Exit(1)
	}

	w, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to get worktree %s\n", err)
	}

	_, err = w.Commit("init", &git.CommitOptions{
		Author:    createSignature(),
		Committer: createSignature(),
	})

	if err != nil {
		fmt.Printf("Failed to commit %s\n", err)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: getSshKeyAuth(),
	})

	if err != nil {
		fmt.Printf("[createAndPushEmptyRepo] Failed to push %s\n", err)
	}
}

func (gs *GitStorage) initRepo() {
	r, err := git.PlainClone(storageDir, false, &git.CloneOptions{
		URL:  CONFIG.Main.Repo,
		Auth: getSshKeyAuth(),
	})

	if repoIsEmptyError(err) {
		fmt.Println("empty repo")
		createAndPushEmptyRepo()
		gs.openRepo()
	} else if repoAlreadyExistsError(err) {
		gs.openRepo()
	} else {
		fmt.Println("new repo")
		gs.repo = r
	}
}

func (gs *GitStorage) commit() {
	w, err := gs.repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to retrieve worktree %s\n", err)
	}

	// git add *
	_, err = w.Add("")
	if err != nil {
		fmt.Printf("Failed to add files %s\n", err)
	}

	// git commit
	_, err = w.Commit("up", &git.CommitOptions{Author: createSignature()})
	if err != nil {
		fmt.Printf("Failed to commit %s\n", err)
	}
}

func (gs *GitStorage) push() {
	err := gs.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{"HEAD:refs/heads/master"},
		Auth: getSshKeyAuth(),
		Prune: true,
	})

	if err != nil {
		fmt.Printf("%s Failed to push %s\n", emojiSave, err)
	}
}

func isUpToDate(err error) bool {
	if err.Error() == "already up-to-date" {
		return true
	}
	return false
}

func (gs *GitStorage) pull() {
	w, err := gs.repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to retrieve worktree %s\n", err)
	}

	err = gs.repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       getSshKeyAuth(),
		Force:      true,
	})

	if isUpToDate(err) == false {
		log.Fatalf("%s Failed to fetch: %s\n", emojiSadFace, err)
	}

	remoteMaster := plumbing.NewRemoteReferenceName("origin", "master")
	err = w.Checkout(&git.CheckoutOptions{Force: true, Branch: remoteMaster})
	if err != nil {
		fmt.Printf("%s Failed to pull: %s\n", emojiSadFace, err)
	}

	fmt.Printf("%s Done\n", emojiCool)
}
