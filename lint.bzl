""""
A test rule that takes one or more `proto_library` targets as input and asserts that
there are no unused import statements in the proto files in the `srcs` attribute
of the target(s).

"""

def _runfile_relative_paths(ctx, proto_import_paths):
    # ProtoInfo provides us with execroot relative paths
    # but we need runfiles relative paths for a test rule
    # and so we need to strip the bin dir prefix wherever necessary.
    # there's probably an elegant way to do this, but this seems to work
    runfiles_import_paths = []
    for imp in proto_import_paths:
        parts = imp.split(ctx.bin_dir.path + "/")

        # this is indeed a bin-dir relative path
        if len(parts) > 1:
            runfiles_import_paths.append(parts[1])
        else:
            # if not, just use as-is
            runfiles_import_paths.append(parts[0])

    return runfiles_import_paths

def _impl(ctx):
    all_proto_sources = []
    all_proto_transitive_sources = []
    all_proto_transitive_imports = []
    for proto in ctx.attr.protos:
        all_proto_sources.extend(proto[ProtoInfo].direct_sources)
        all_proto_transitive_sources.append(proto[ProtoInfo].transitive_sources)
        all_proto_transitive_imports.extend(proto[ProtoInfo].transitive_proto_path.to_list())

    test_script = ctx.actions.declare_file("{}_extended.sh".format(ctx.label.name))

    all_proto_transitive_imports = _runfile_relative_paths(ctx, all_proto_transitive_imports)

    ctx.actions.expand_template(
        template = ctx.file._test_template,
        output = test_script,
        is_executable = True,
        substitutions = {
            "%%unused_imports%%": ctx.executable._unused_imports_tool.short_path,
            "%%proto_files%%": ",".join([p.path for p in all_proto_sources]),
            "%%import_paths%%": ",".join(all_proto_transitive_imports),
        },
    )

    runfiles = ctx.runfiles(files = [ctx.executable._unused_imports_tool] + ctx.files._runfiles_bash_lib)

    # We need all the proto sources, and their transitive imports
    transitive_runfiles = ctx.runfiles(
        files = all_proto_sources,
        transitive_files = depset(
            direct = [],
            order = "default",
            transitive = all_proto_transitive_sources,
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
        "protos": attr.label_list(
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
