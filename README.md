CryptoFungi Project

# 크립토풍기 프로젝트 실행 방법

## 1. BlockChain Network 실행

./startNetwork.sh

## 2. 체인코드 배포
### 버섯공장(fungusfactory)
./deployFungiCC.sh 1 1

### 먹이공장(feedfactory)
./deployFeedCC.sh 1 1

## 3. 체인코드 테스트 ( Initalize를 위해 반드시 수행필요 )
### 먹이공장(feedfactory)
./testFeedCC.sh

### 버섯공장(fungusfactory)
./testFungiCC.sh


### 버섯거래 (feedfactory)
./testFungiCC2.sh
* 카우치DB에서 생성된 버섯의 Owner id를 조회하여 보내는사람과 받는사람 부분의 args를 변경하여 수행
* 카우치DB 확인 url : http://localhost:5984/_utils/

## 4. Application 실행
./startApplication.sh

## 5. Web Client 접속
웹브라우저 : localhost:3000

## END : 네트워크 종료
./downNetwork.sh