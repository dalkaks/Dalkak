package application

import (
	"dalkak/internal/infrastructure/eventbus"
	userdto "dalkak/pkg/dto/user"
	cryptoutil "dalkak/pkg/utils/crypto"
	jwtutil "dalkak/pkg/utils/jwt"
	responseutil "dalkak/pkg/utils/response"
	timeutil "dalkak/pkg/utils/time"
)

func (app *ApplicationImpl) RegisterUserEventListeners() {
	app.EventManager.Subscribe("post.user.auth", app.handleAuthAndSignUp)
	app.EventManager.Subscribe("post.user.refresh", app.handleReissueAccessToken)
}

func (app *ApplicationImpl) handleAuthAndSignUp(event eventbus.Event) {
	payload, ok := event.Payload.(*userdto.AuthAndSignUpRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

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
	accessToken, refreshToken, err := jwtutil.GenerateAuthToken(app.Keymanager, &jwtutil.GenerateTokenDto{
		WalletAddress: payload.WalletAddress,
		NowTime:       timeutil.GetTimestamp(),
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := userdto.NewAuthAndSignUpResponse(accessToken.Token, accessToken.TokenTTL, refreshToken.Token, refreshToken.TokenTTL)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}

func (app *ApplicationImpl) handleReissueAccessToken(event eventbus.Event) {
	payload, ok := event.Payload.(*userdto.ReissueAccessTokenRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 토큰 발급
	accessToken, err := jwtutil.GenerateAccessToken(app.Keymanager, &jwtutil.GenerateTokenDto{
		WalletAddress: payload.WalletAddress,
		NowTime:       timeutil.GetTimestamp(),
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := userdto.NewReissueAccessTokenResponse(accessToken.Token, accessToken.TokenTTL)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
