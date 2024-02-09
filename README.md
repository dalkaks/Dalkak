# Dalkak

# 서버

## 사전 설치

- gnu make (gmake)
- go 1.21.x
- aws-cli

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
 - 현재 커밋을 태그로 배포 (덮어쓰기 시 아래 명령어에 -f 추가)
 - 명령어 git tag test-v0.0.0
 - 명령어 git push origin test-v0.0.0
 - 이후 슬랙으로 결과 전송
 - ecs 정지를 위해 명령어 ./script/dev-server-switch.sh
 