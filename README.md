# HOW TO BUILD

[MANUAL-PDF](./manual.pdf)

## SERVER

- 인증서 교체가 필요하다면 아래 optional 수행 필요
- docker 설치 필요

```bash
# 빌드
docker build . -t webmote

# 실행 
docker run -it --name webmote -p 443:8000 webmote
```

### (*optional) SSL 용 인증서 생성

- 본인이 보유하고 있는 인증서가 존재한다면 `tls.key`, `tls.cert` 교체 필요  
- 없다면 새로 생성 가능
- `CN=YOURDOMAIN_OR_IP` 수정 필요

```bash
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.cert -sha256 -days 3650 -nodes -subj "/O=WEBMOTE/OU=WEBMOTE/CN=YOURDOMAIN_OR_IP"
```

## GAME

- `godot-engine` 설치 필요
- 자체 서버를 활용한다면 도메인 재정의 필요
- INVALID 한 인증서를 활용한다면 인증서 신뢰 설정 필요

1. [GODOT REPO](https://github.com/godotengine/godot/releases/) 접속 후 설치
2. [GAME REPO](https://github.com/osswebmote/fpsGame) 클론 혹은 다운로드
3. GODOT 에서 프로젝트 OPEN 및 빌드 가능

### (*optional) 도메인 재정의

아래 스크립트의 `YOUR_OWN_DOMAIN_OR_IP` 재정의 필요

- 맥을 사용하는 경우

```bash
sed 's/demo.mansuiki.com/YOUR_OWN_DOMAIN_OR_IP/g' ./script/network/global_data.gd > ./script/network/global_data.gd.new
mv ./script/network/global_data.gd.new ./script/network/global_data.gd
```

- 리눅스를 사용하는 경우

```bash
sed -i 's/demo.mansuiki.com/YOUR_OWN_DOMAIN_OR_IP/g' ./script/network/global_data.gd
```

### (*optional) 인증서 신뢰 설정

1. 위 `(*optional) SSL 용 인증서 생성` 에서 생성한 tls.cert 파일의 이름을 `tls.crt` 로 변경후 프로젝트 폴더 내부에 복사
2. Godot의 `Project Settings` -> 설정 필터에 `tls` 입력
3. `네트워크 > TLS > 인증서 번들 오버라이드` 에서 `tls.crt` 선택
4. 게임 빌드 후 실행
