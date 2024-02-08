# Dalkak

# 서버

## 사전 설치

- gnu make 설치
   - [Windows](https://jstar0525.tistory.com/264)
- [go 1.21.x](https://go.dev/dl/)
- [aws-cli](https://docs.aws.amazon.com/ko_kr/cli/latest/userguide/getting-started-install.html)

## aws-cli 설정

 - 명령어 aws configure --profile dalkak
 - region: ap-northeast-2
 - output: json

## 서버 실행

 - make -C server run-local
 (또는 server 폴더에서 make run-local)


# 클라이언트

# 배포

## develop 배포
 - 현재 커밋의 경우
 - 명령어 git tag test-v0.0.0
 - 명령어 git push origin test-v0.0.0
 (덮어쓰기 시 -f 옵션 추가)
 - 이후 슬랙으로 결과 전송
 - ecs 정지를 위해 명령어 ./script/stop-dev-server.sh
 
