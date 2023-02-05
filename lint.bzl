""""
A test rule that takes a `proto_library` target as input and asserts that
there are no unused import statements in the proto files in the `srcs` attribute
of the target.

"""

def _impl(ctx):
    proto_info = ctx.attr.proto[ProtoInfo]

    test_script = ctx.actions.declare_file("{}_unused_imports.sh".format(ctx.label.name))
    package = ctx.label.package
    transitive_import_paths = proto_info.transitive_proto_path.to_list()

    ctx.actions.expand_template(
        template = ctx.file._test_template,
        output = test_script,
        is_executable = True,
        substitutions = {
            "%%unused_imports%%": ctx.executable._unused_imports_tool.short_path,
            "%%protos_dir%%": package,
            "%%import_paths%%": ",".join(transitive_import_paths),
        },
    )

    runfiles = ctx.runfiles(files = [ctx.executable._unused_imports_tool] + ctx.files._runfiles_bash_lib)


    # We need all the proto sources, and their transitive imports
    transitive_runfiles = ctx.runfiles(
        files = proto_info.direct_sources,
        transitive_files = depset(
            direct = [],
            order = "default",
            transitive = [proto_info.transitive_sources],
        ),
    )

    runfiles = runfiles.merge(transitive_runfiles)

    return [
        DefaultInfo(
            executable = test_script,
            runfiles = runfiles,
        ),
    ]

proto_unused_imports_test = rule(
    implementation = _impl,
    test = True,
    attrs = {
        "proto": attr.label(
            providers = [ProtoInfo],
            doc = "A proto_library which is to be linted for unused imports.",
        ),
        "_unused_imports_tool": attr.label(
            default = Label("//internal/unused_imports"),
            executable = True,
            cfg = "exec",
        ),
        "_runfiles_bash_lib": attr.label(
            default = Label("@bazel_tools//tools/bash/runfiles"),
        ),
        "_test_template": attr.label(
            default = Label("//internal:lint_test.sh.template"),
            allow_single_file = True,
        ),
    },
)
