package usecase

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"go.uber.org/zap"
	gitHttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)



func (u *Usecase) GetOrdinalTemplate() (*os.File, error)  {
	_, err :=  helpers.FileExists("project.zip")
	if err == nil {
		f, err  := os.Open("project.zip")
		if err == nil {
			return f, nil
		}
	}

	metaJson := structure.OrdinalCollectionMeta{
		Name: "name",
		InscriptionIcon: "inscription_icon",
		Supply: "1",
		Slug: "slug",
		Description: "description",
		TwitterLink: "twitter_link",
		DiscordLink: "discord_link",
		WebsiteLink: "website_link",
		WalletAddress: "wallet_address",
		Royalty: "0",
	}
	metaName := `meta.json`
	inscriptionName := `inscriptions.json`
	helpers.CreateFile(metaName, metaJson)
	iMeta :=  make(map[string]string)
	iMeta["name"] = "name"
	inscriptionsJson := structure.OrdinalInscriptionMeta{
		ID: "inscription_id",
		Meta: iMeta,
	}
	inscriptionList := []structure.OrdinalInscriptionMeta{}
	inscriptionList = append(inscriptionList, inscriptionsJson)
	helpers.CreateFile(inscriptionName, inscriptionList)
	
	zipFile, err := helpers.ZipFiles("project", metaName, inscriptionName)
	if err != nil {
		return nil, err
	}
	return zipFile, nil
}

func (u *Usecase) UploadOrdinalTemplate(r *http.Request) (*os.File, error)  {
	
	projectName := ""
	zippedBytes,folderName, err := u.uploadAndReadFile(r, "file")
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(zippedBytes)
    zipReader, err := zip.NewReader(reader, int64(len(zippedBytes)))
    if err != nil {
        return nil, err
    }

	metaCollection  := &structure.OrdinalCollectionMeta{}
	metaInscriptions  := []structure.OrdinalInscriptionMeta{}
	for _, file := range zipReader.File {
		fc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fc.Close()

		content, err := ioutil.ReadAll(fc)
		if err != nil {
			return nil, err
		}

		
		if strings.Index(file.Name, "meta") != -1 {
			err = json.Unmarshal(content, metaCollection)
			if err != nil {
				return nil, err
			}
		}
		
		if strings.Index(file.Name, "inscriptions") != -1 {
			err = json.Unmarshal(content, &metaInscriptions)
			if err != nil {
				return nil, err
			}
		}
	}

	err = metaCollection.Verify()
	if err != nil {
		return nil, err
	}

	if len(metaInscriptions) <= 0 {
		err := errors.New("Plase upload inscription.json")
		return nil, err
	}

	inscription := metaInscriptions[0].ID
	for i, inscription := range metaInscriptions {
		err := inscription.Verify()
		if err != nil {
			err = fmt.Errorf("Inscription %d failure with error %v", i+1, err)
			return nil, err
		}

		token, err := u.Repo.FindTokenByTokenID(inscription.ID)
		if err == nil && token != nil {
			err = fmt.Errorf("Inscription %s at line %d is existed",inscription.ID, i+1)
			return nil, err
		}

	}

	if metaCollection.InscriptionIcon != inscription {
		metaCollection.InscriptionIcon = inscription
	}

	
	err = u.pushToGithub(*folderName, *metaCollection, metaInscriptions)
	if err != nil {
		return nil, err
	}
	logger.AtLog.Infof("collection %s ",projectName)
	return nil, err
}

func (u Usecase) uploadAndReadFile(r *http.Request, fileName string) ([]byte, *string, error) {
	_, handler, err := r.FormFile(fileName)
	key := "btc-projects/ordinal-inscriptions"
	gf := googlecloud.GcsFile{
		FileHeader: handler,
		Path:       &key,
	}

	
	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", key), err.Error())
		return nil, nil, err
	}

	content, err := u.GCS.ReadFile(uploaded.Name)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", key), err.Error())
		return nil,nil, err
	}

	folerNames := strings.Split(uploaded.Name, "/")
	folerName := folerNames[len(folerNames) - 1]
	
	return content, &folerName, nil
}

func (u Usecase) pushToGithub(folderName string, metaCollection structure.OrdinalCollectionMeta, inscription []structure.OrdinalInscriptionMeta) error   {
	source := "https://github.com/generative-xyz/ordinals-collections.git"
	uuid := uuid.New().String()
	
	//timeUtc := time.Now().UTC().Nanosecond()
	//folderName = fmt.Sprintf("%s-%d",folderName, timeUtc)
	folder_path := fmt.Sprintf("/tmp/ordinals-collection-%s", uuid)
	collectionPath := fmt.Sprintf("%s/collections/%s",folder_path, folderName)
	repo, err := git.PlainClone(folder_path, false, &git.CloneOptions{
		URL:      source,
		//Progress: os.Stdout,
	})

	if err != nil {
		return  err
	}
	collectionFoldersPath := fmt.Sprintf("%s/collections", folder_path)
	collectionFolders, err := ioutil.ReadDir(collectionFoldersPath)
	if err != nil {
		return  err
	}

	//create new branch and checkout it
	w, err := repo.Worktree()
	if err != nil {
		return  err
	}

	headRef, err := repo.Head()
	if err != nil {
		return  err
	}

	branch := os.Getenv("GITHUB_ORDINAL_COLLECTOR_BRANCH")
	ref := plumbing.NewHashReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)), headRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		return  err
	}

	spew.Dump(collectionPath)
	branchName :=  plumbing.NewBranchReferenceName(branch)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: branchName,
	})
	if err != nil {
		return  err
	}

	err = w.Pull(&git.PullOptions{
			RemoteName: `origin`,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
		},
	)
	
	//_ = ref
	_ = collectionFolders
	spew.Dump(collectionPath)
	//spew.Dump(repo)
	//spew.Dump(ref)

	//create files
	// for _, f := range collectionFolders {
	// 	spew.Dump(f.Name())
	// 	_ = f
	// }

	err = os.Mkdir(collectionPath, os.ModePerm)
	if err != nil {
		return  err
	}

	err = helpers.CreateFile(fmt.Sprintf("%s/meta.json",collectionPath), metaCollection)
	if err != nil {
		return  err
	}

	err = helpers.CreateFile(fmt.Sprintf("%s/inscriptions.json",collectionPath), metaCollection)
	if err != nil {
		return  err
	}

	//add . and commmit them
	_, err = w.Add(".")
	if err != nil {
		return  err
	}

	status, err := w.Status()
	if err != nil {
		return  err
	}

	// commit and push
	commit, err := w.Commit(fmt.Sprintf("Collection %s is created",folderName), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "generative",
			Email: "dev@generative.xyz",
			When:  time.Now(),
		},
	})
	if err != nil {
		return  err
	}

	err = repo.Push(&git.PushOptions{
		//RemoteName: os.Getenv("GITHUB_ORDINAL_COLLECTOR_REPO"),
		Auth: &gitHttp.BasicAuth{
			Username: os.Getenv("GITHUB_ORDINAL_COLLECTOR_OWNER"), 
			Password: os.Getenv("GITHUB_ACCESSTOKEN"),
		},
		//Force: true,
	})
	if err != nil {
		return  err
	}

	logger.AtLog.Info(zap.String("collectionPath", collectionPath), zap.String("branch", branch), zap.Any("commit",commit), zap.Any("status",status))
	// spew.Dump(commit)
	// spew.Dump(status)
	// spew.Dump(commit)
	return  nil
}