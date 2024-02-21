package application

import (
	"dalkak/internal/infrastructure/eventbus"
	mediadto "dalkak/pkg/dto/media"
	userdto "dalkak/pkg/dto/user"
	cryptoutil "dalkak/pkg/utils/crypto"
	jwtutil "dalkak/pkg/utils/jwt"
	timeutil "dalkak/pkg/utils/time"
)

func (app *ApplicationImpl) RegisterUserEventListeners() {
	app.EventManager.Subscribe("post.user.auth", app.handleAuthAndSignUp)
	app.EventManager.Subscribe("post.user.refresh", app.handleReissueAccessToken)
	app.EventManager.Subscribe("post.user.media.presigned", app.handleCreateTempMedia)
}

func (app *ApplicationImpl) handleAuthAndSignUp(event eventbus.Event) {
	// ok
	payload := event.Payload.(*userdto.AuthAndSignUpRequest)

	// 시그니처 검증
	err := cryptoutil.VerifyMetaMaskSignature(&cryptoutil.MetaMaskSignature{
		WalletAddress: payload.WalletAddress,
		Signature:     payload.Signature,
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 유저 조회 및 생성
	checkAndCreateUserDto := userdto.NewCheckAndCreateUserDto(payload.WalletAddress)
	newUser, err := app.UserDomain.CheckAndCreateUser(checkAndCreateUserDto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 유저 저장
	if newUser != nil {
		err := app.Database.CreateUser(newUser)
		if err != nil {
			app.SendResponse(event.ResponseChan, nil, err)
			return
		}
	}

	// 토큰 발급
	accessToken, refreshToken, err := jwtutil.GenerateAuthToken(app.AppConfig.Domain, app.Keymanager, &jwtutil.GenerateTokenDto{
		WalletAddress: payload.WalletAddress,
		NowTime:       timeutil.GetTimestamp(),
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := userdto.NewAuthAndSignUpResponse(accessToken.Token, accessToken.TokenTTL, refreshToken.Token, refreshToken.TokenTTL)
	app.SendResponse(event.ResponseChan, result, nil)
}

func (app *ApplicationImpl) handleReissueAccessToken(event eventbus.Event) {
	// ok
	payload := event.Payload.(*userdto.ReissueAccessTokenRequest)

	// 토큰 발급
	accessToken, err := jwtutil.GenerateAccessToken(app.AppConfig.Domain, app.Keymanager, &jwtutil.GenerateTokenDto{
		WalletAddress: payload.WalletAddress,
		NowTime:       timeutil.GetTimestamp(),
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := userdto.NewReissueAccessTokenResponse(accessToken.Token, accessToken.TokenTTL)
	app.SendResponse(event.ResponseChan, result, nil)
}

func (app *ApplicationImpl) handleCreateTempMedia(event eventbus.Event) {
	// ok
	userInfo := event.UserInfo
	payload := event.Payload.(*userdto.CreateTempMediaRequest)

	// 미디어 생성
	dto := mediadto.NewCreateTempMediaDto(userInfo, payload.MediaType, payload.Ext, payload.Prefix)
	newMedia, err := app.MediaDomain.CreateMediaTemp(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 미디어 저장
	err = app.Database.CreateUserMediaTemp(userInfo.GetUserId(), newMedia)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := userdto.NewUserCreateMediaResponse(newMedia.MediaEntity.Id, newMedia.MediaTempUrl.AccessUrl, *newMedia.MediaTempUrl.UploadUrl)
	app.SendResponse(event.ResponseChan, result, nil)
}
