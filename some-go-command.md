
```
빌드 결과 캐시 삭제
$ go clean -cache

테스트 결과 캐시 삭제
$ go clean -testcache

mod 패키지 삭제
go clean -modcache

테스트 태그
go test --tags=integration ./...
```