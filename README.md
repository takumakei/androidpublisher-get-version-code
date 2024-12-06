# androidpublisher-get-version-code
Get latest version code from Google Play Developer API

- [VersionCode](https://pkg.go.dev/google.golang.org/api@v0.210.0/androidpublisher/v3#TrackRelease.VersionCodes)
- [TracksListResponse](https://pkg.go.dev/google.golang.org/api@v0.210.0/androidpublisher/v3#TracksListResponse)

## Install

You can download pre-built binaries from [Releases](https://github.com/takumakei/androidpublisher-get-version-code/releases).

Or, you can build and install using Go like below:

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

  -jmespath-expr JMESPATH_EXPR
        JMESPath expression applying to the result.
        If not given, the value will be checked from the environment variable JMESPATH_EXPR.

  -output-style OUTPUT_STYLE
        Output style selector:
          - highest
            The highest build number across all tracks
          - production, beta, alpha, internal
            Select the track
          - response
            The JSON of TracksListResponse
        If not given, the value will be checked from the environment variable OUTPUT_STYLE.
        If none is specified, `highest` is used.

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

## Credentials

See [Collect your Google credentials](https://docs.fastlane.tools/getting-started/android/setup/#collect-your-google-credentials).

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
