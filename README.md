## proto-lint-unused-imports
This repository exposes a simple Bazel rule that attempts to solve (or workaround) the issue described in https://github.com/bufbuild/rules_buf/issues/32.

Currently the buf plugin that runs the `buf lint` tests doesn't fail on unused imports for reasons described in the issue. This means the Bazel rulesets that do exposes the buf lint rules do not fail on unused imports.

## Approach
This ends up utilizing much of the buf protoc engine at https://github.com/bufbuild/protocompile but combines it with a custom test rule `proto_unused_imports_test` that explicitly operates on the proto sources described by a `proto_library` target. A custom linter written in `Go` utilizing the Buf APIs then acts as a test runner and then validates that the unused imports do end up failing the targets.

## Pre-requisites
You need [Bazel](https://bazel.build/install) installed on your machine.


## Examples
There is currently a simple example in `examples/` directory. In `bar.proto` there are two unused imports,
```
import "google/protobuf/any.proto";
import "google/protobuf/descriptor.proto";
```

If you run `bazel test //examples/simple:simple_proto_test`, you'll notice that the test will fail with the unused imports.

```
INFO: Analyzed target //examples/simple:simple_proto_test (0 packages loaded, 0 targets configured).
INFO: Found 1 test target...
FAIL: //examples/simple:simple_proto_test (see /private/var/tmp/_bazel_smocherla/8f6414120d5c975fe9a914425066b8a5/execroot/proto-lint-unused-imports/bazel-out/darwin_arm64-fastbuild/testlogs/examples/simple/simple_proto_test/test.log)
INFO: From Testing //examples/simple:simple_proto_test:
==================== Test output for //examples/simple:simple_proto_test:
Found 2 protos
Detected unused import in examples/simple/bar.proto at line 5
import "google/protobuf/any.proto" not used
Detected unused import in examples/simple/bar.proto at line 6
import "google/protobuf/descriptor.proto" not used
```


## Improvements
This is currently tested with well-known types and their imports but should also work with workspace-relative imports as all transitive import paths are passed through to Buf's compiler engine.

More examples can be added to evaluate various edge cases.

This lint can serve as an example of "extending" Bufs's API to write your own custom linting rules on any aspect of a proto AST.