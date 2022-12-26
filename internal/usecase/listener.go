package usecase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
)

type projectChan struct {
	Data *entity.Projects
	Err error
}

type projectDetailChan struct {
	Data *structure.ProjectDetail
	Err error
}

func (u Usecase) UpdateProjectWithListener(chainLog types.Log) {
	txnHash := chainLog.TxHash.String()
	span, log := u.StartSpanWithoutRoot("UpdateProjectWithListener.GetProjectDetail")
	defer u.Tracer.FinishSpan(span, log )
	log.SetTag("transaction_hash", txnHash)
	log.SetData("chainLog", chainLog)
	topics := chainLog.Topics

	
	
	tokenIDStr :=  helpers.HexaNumberToInteger(topics[3].String())
	tokenID, _ := strconv.Atoi(tokenIDStr)
	tokenIDStr = fmt.Sprintf("%d",tokenID)
	contractAddr := strings.ToLower(chainLog.Address.String()) 

	u.UpdateProjectFromChain(contractAddr, tokenIDStr)
}

func (u Usecase) UpdateProjectFromChain(contractAddr string, tokenIDStr string) (*entity.Projects, error) {
	span, log := u.StartSpanWithoutRoot("UpdateProjectWithListener.GetProjectDetail")
	defer u.Tracer.FinishSpan(span, log )

	pChan := make(chan projectChan, 1)
	pDChan := make(chan projectDetailChan, 1)
	
	log.SetData("contractAddr", contractAddr)
	log.SetData("tokenIDStr", tokenIDStr)

	go func(rootSpan opentracing.Span, pChan chan projectChan, contractAddr string, tokenIDStr string) {
		span, log := u.StartSpan("GetProjectDetail.Project", rootSpan)
		defer u.Tracer.FinishSpan(span, log )

		log.SetTag("contractAddr",contractAddr)
		log.SetTag("tokenIDStr",tokenIDStr)

		project := &entity.Projects{}
		var err error

		defer func  ()  {
			pChan <- projectChan {
				Data: project,
				Err:  err,
			}
		}()

		project, err = u.Repo.FindProjectBy(contractAddr, tokenIDStr)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				project = &entity.Projects{}
				project.ContractAddress = contractAddr 
				project.TokenID = tokenIDStr 

				err = u.Repo.CreateProject(project)
				if err != nil {
					log.Error("u.Repo.CreateProject", err.Error(), err)
					return
				}

		}else{
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return
		}
	}

	}(span, pChan, contractAddr, tokenIDStr)

	go func(rootSpan opentracing.Span, pDChan chan projectDetailChan, contractAddr string, tokenIDStr string) {
		span, log := u.StartSpan("GetProjectDetail.Project", rootSpan)
		defer u.Tracer.FinishSpan(span, log )
		
		projectDetail := &structure.ProjectDetail{}
		var err error

		defer func  ()  {
			pDChan <- projectDetailChan {
				Data: projectDetail,
				Err:  err,
			}
		}()

		projectDetail, err = u.getProjectDetailFromChain(span, structure.GetProjectDetailMessageReq{
			ContractAddress:  contractAddr,
			ProjectID:  tokenIDStr,
		})
		if err != nil {
			log.Error(" u.GetProjectDetail", err.Error(), err)
			return
		}
	

	}(span, pDChan, contractAddr, tokenIDStr)
	
	projectFChan := <- pChan
	projectDetailFChan := <- pDChan

	err := projectFChan.Err 
	if err != nil {
		log.Error("projectFChan.Err ", err.Error(), err)
		return nil, err
	}

	err = projectDetailFChan.Err
	if err != nil {
		log.Error("projectDetailFChan.Err ", err.Error(), err)
		return nil, err
	}

	project := projectFChan.Data
	projectDetail := projectDetailFChan.Data

	log.SetData("project", project)
	log.SetData("projectDetail", projectDetail)

	project.Name = projectDetail.ProjectDetail.Name
	project.CreatorName = projectDetail.ProjectDetail.Creator
	project.CreatorAddrr = projectDetail.ProjectDetail.CreatorAddr.String()
	project.Description = projectDetail.ProjectDetail.Desc
	project.Scripts= projectDetail.ProjectDetail.Scripts
	project.ThirdPartyScripts= projectDetail.ProjectDetail.ScriptType
	project.Styles= projectDetail.ProjectDetail.Styles
	project.GenNFTAddr= projectDetail.ProjectDetail.GenNFTAddr.String()
	//project.Hash = txnHash
	project.MintPrice = projectDetail.ProjectDetail.MintPrice.String()
	project.MaxSupply = int(projectDetail.ProjectDetail.MaxSupply.Int64())
	project.LimitSupply = int(projectDetail.ProjectDetail.Limit.Int64())
	project.MintTokenAddress = string(projectDetail.ProjectDetail.MintPriceAddr.String())
	project.License = projectDetail.ProjectDetail.License
	project.Status = projectDetail.Status
	project.SocialWeb = projectDetail.ProjectDetail.Social.Web
	project.SocialTwitter = projectDetail.ProjectDetail.Social.Twitter
	project.SocialDiscord = projectDetail.ProjectDetail.Social.Discord
	project.SocialMedium = projectDetail.ProjectDetail.Social.Medium
	project.SocialInstagram = projectDetail.ProjectDetail.Social.Instagram
	project.Thumbnail = projectDetail.ProjectDetail.Image
	project.NftTokenUri = projectDetail.NftTokenUri
	project.IsSynced = true

	updated, err := u.Repo.UpdateProject(project.UUID, project)
	if err != nil {
		log.Error(" u.UpdateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("updated",updated)
	return  project, nil
}

