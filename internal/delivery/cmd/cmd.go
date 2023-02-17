package cmd

import (
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"
)


type CMDHandler struct {
	Logger   logger.Ilogger
	Tracer   tracer.ITracer
	Cache    redis.IRedisCache
	Usecase    usecase.Usecase
}

func NewCMDHandler(global *global.Global, uc usecase.Usecase) *CMDHandler {
	return &CMDHandler{
		Logger: global.Logger,
		Tracer: global.Tracer,
		Cache: global.Cache,
		Usecase: uc,
	}
}


func (h CMDHandler) StartServer() {	
	//h.Usecase.LoadImage()
	// h.Usecase.ChangeProjectCreatorProfile("1000103","0xDA1958529ACCed8834FEf1D0e48a8cebD618f159" )
	// h.Usecase.ChangeProjectCreatorProfile("1000104","0xDA1958529ACCed8834FEf1D0e48a8cebD618f159" )
	// h.Usecase.ChangeProjectCreatorProfile("1000105","0xDA1958529ACCed8834FEf1D0e48a8cebD618f159" )
	// h.Usecase.ChangeProjectCreatorProfile("1000107","0xDA1958529ACCed8834FEf1D0e48a8cebD618f159" )
	h.Usecase.ChangeProjectCreatorProfile("1000113","0xF9BE59246C5097D80c40c112Af9678b2E06eA939" )
	// h.Usecase.ChangePrice("1000101", "0.0049")
	// h.Usecase.ChangePrice("1000112", "0.0069")
	// h.Usecase.DeleteProjectID("1000133")
	// h.Usecase.ChangeRoyalty("1000140", 5*100)
	// h.Usecase.ChangeRoyalty("1000140", 5*100)
	//h.Usecase.UpdateProfileProfile("1000067")
	//h.Usecase.Update1M02Collections("1000122")
}
