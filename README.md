# androidpublisher-get-version-code
Get latest version code from Google Play Developer API

- [VersionCode](https://pkg.go.dev/google.golang.org/api@v0.210.0/androidpublisher/v3#TrackRelease.VersionCodes)
- [TracksListResponse](https://pkg.go.dev/google.golang.org/api@v0.210.0/androidpublisher/v3#TracksListResponse)

## Install

```
go install github.com/takumakei/androidpublisher-get-version-code@latest
```

## Usage

```
Usage of androidpublisher-get-version-code:
  -credentials CREDENTIALS
        Gcloud service account credentials with the JSON key type to access Google Play Developer API.
        If not given, the value will be checked from the environment variable CREDENTIALS.
        Alternatively to entering CREDENTIALS in plaintext, it may also be specified using the
        "@env:" prefix followed by an environment variable name, or the
        "@file:" prefix followed by a path to the file containing the value.

  -output-style OUTPUT_STYLE
        Output style selector:
          - highest
            The highest build number across all tracks
          - production, beta, alpha, internal
            Select the track
          - response
            The JSON of TracksListResponse
        If not given, the value will be checked from the environment variable OUTPUT_STYLE.

  -package-name PACKAGE_NAME
        Package name of the app in Google Play Console. For example com.example.app
        If not given, the value will be checked from the environment variable PACKAGE_NAME.

  -time-limit TIME_LIMIT
        Specifies the duration for how long to wait for a response from the API call.
        If not given, the value will be checked from the environment variable TIME_LIMIT.
         (default 30s)
  -version
        print version information
```

## Example

```
$ androidpublisher-get-version-code --package-name com.example.app --output-style internal
{
  "track": "internal",
  "name": "1.0.0",
  "code": 1
}
$
```

## Related work

https://github.com/codemagic-ci-cd/cli-tools/blob/master/docs/google-play/get-latest-build-number.md
