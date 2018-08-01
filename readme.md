사용법

LineChanger configjson 빌드단계 (프로젝트경로)

json 포맷 - 첨부된 sample.json을 참고
* 배열
 * path : 고칠 파일 프로젝트경로+path
 * description : 그냥 주석용
 * encoding : 현재 euckr, utf16bom, utf8 3가지만 지원
 * change : 배열
  * find : 찾을 패턴 : 앞내용[[tochange]](뒷내용)
  * description : 그냥 주석용
  * changeto : 맵
   * buildStage : 바꿀 내용

예를 들어 

* "find": "sdkKey = TEXT (\"[[tochange]]\");\t// sdkKey",
* 이런 문장을 입력하면
* sdkKey = TEXT ("어쩌구저쩌구");	// sdkKey
* 이런 문장이 걸리고
* sdkKey = TEXT ("asldjaslkjqawlfdjwqlf");	// sdkKey
* 이런 식으로 바뀌어서 저장한다
