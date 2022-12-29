package usecase

import (
	"github.com/opentracing/opentracing-go"
	nftStructure "rederinghub.io/external/nfts"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetNftTransactions(rootSpan opentracing.Span, req structure.GetNftTransactionsReq) (*nftStructure.CovalentGetNftTransactionResponse, error) {
	span, log := u.StartSpan("GetNftTransactions", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req);
	resp, err := u.CovalentNft.GetNftTransactions(nftStructure.CovalentNftTransactionFilter{
		Chain : req.Chain,
		ContractAddress: req.ContractAddress,
		TokenID: req.TokenID,
	})
	return resp, err
}
