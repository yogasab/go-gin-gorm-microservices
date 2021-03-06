package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/service"
	"net/http"
	"strconv"
)

type CampaignHandler interface {
	Index(ctx *gin.Context)
	Add(ctx *gin.Context)
	Store(ctx *gin.Context)
	UploadImage(ctx *gin.Context)
	StoreImage(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Update(ctx *gin.Context)
	Show(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService service.CampaignService
	userService     service.UserService
}

func NewCampaignHandler(campaignService service.CampaignService, userService service.UserService) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
		userService:     userService,
	}
}

func (h *campaignHandler) Index(ctx *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	ctx.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) Add(ctx *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	request := dto.FormStoreCampaignRequest{}
	request.Users = users
	ctx.HTML(http.StatusOK, "add_campaign.html", request)
}

func (h *campaignHandler) Store(ctx *gin.Context) {
	var request dto.FormStoreCampaignRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		users, e := h.userService.GetAllUsers()
		if e != nil {
			ctx.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		request.Error = err
		request.Users = users
		ctx.HTML(http.StatusInternalServerError, "error.html", request)
		return
	}
	user, err := h.userService.GetUserByID(request.UserID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	formRequest := dto.CreateCampaignRequest{}
	formRequest.Name = request.Name
	formRequest.ShortDescription = request.ShortDescription
	formRequest.Description = request.Description
	formRequest.Perks = request.Perks
	formRequest.GoalAmount = request.GoalAmount
	formRequest.User = user
	_, err = h.campaignService.CreateCampaign(formRequest)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
	}
	ctx.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) UploadImage(ctx *gin.Context) {
	id := ctx.Param("id")
	campaignId, _ := strconv.Atoi(id)
	ctx.HTML(http.StatusOK, "image_campaign.html", gin.H{"ID": campaignId})
}

func (h *campaignHandler) StoreImage(ctx *gin.Context) {
	//	Take the file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.HTML(http.StatusBadGateway, "error.html", nil)
		return
	}
	//	Take the campaign param id
	id := ctx.Param("id")
	campaignID, _ := strconv.Atoi(id)
	//	Query to get campaign by id
	campaign, err := h.campaignService.GetCampaignByID(dto.CampaignGetRequestID{ID: campaignID})
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	//	Save uploaded file
	userID := campaign.UserId
	destination := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = ctx.SaveUploadedFile(file, destination)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	request := dto.CreateCampaignImageRequest{}
	request.CampaignID = campaignID
	request.IsPrimary = true
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	request.User = user
	//	Take the save campaign image service
	_, err = h.campaignService.CreateCampaignImage(request, destination)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	ctx.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(ctx *gin.Context) {
	id := ctx.Param("id")
	campaignID, _ := strconv.Atoi(id)
	campaign, err := h.campaignService.GetCampaignByID(dto.CampaignGetRequestID{ID: campaignID})
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	request := dto.FormUpdateCampaignRequest{}
	request.ID = campaignID
	request.Name = campaign.Name
	request.ShortDescription = campaign.ShortDescription
	request.Description = campaign.Description
	request.Perks = campaign.Perks
	request.GoalAmount = campaign.GoalAmount
	ctx.HTML(http.StatusOK, "edit_campaign.html", request)
}

func (h *campaignHandler) Update(ctx *gin.Context) {
	// 	Take the campaign id param
	id := ctx.Param("id")
	campaignID, _ := strconv.Atoi(id)
	//	Declare the update request and bind it
	request := dto.FormUpdateCampaignRequest{}
	err := ctx.ShouldBind(&request)
	if err != nil {
		fmt.Println(err.Error())
		request.ID = campaignID
		request.Error = err
		ctx.HTML(http.StatusInternalServerError, "error.html", request)
		return
	}
	//	Query campaign by id
	campaign, err := h.campaignService.GetCampaignByID(dto.CampaignGetRequestID{ID: campaignID})
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userID := campaign.UserId
	// Query user who owns campaign by user id
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	// Mapping the request from update campaign
	formUpdateCampaign := dto.CreateCampaignRequest{}
	formUpdateCampaign.Name = request.Name
	formUpdateCampaign.Perks = request.Perks
	formUpdateCampaign.GoalAmount = request.GoalAmount
	formUpdateCampaign.Description = request.Description
	formUpdateCampaign.ShortDescription = request.ShortDescription
	formUpdateCampaign.User = user
	_, err = h.campaignService.UpdateCampaign(dto.CampaignGetRequestID{ID: campaignID}, formUpdateCampaign)
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	ctx.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	campaignID, _ := strconv.Atoi(id)

	request := dto.CampaignGetRequestID{}
	request.ID = campaignID

	campaign, err := h.campaignService.GetCampaignByID(request)
	if err != nil {
		ctx.Redirect(http.StatusInternalServerError, "/campaigns")
		return
	}
	ctx.HTML(http.StatusOK, "show_campaign.html", campaign)
}
