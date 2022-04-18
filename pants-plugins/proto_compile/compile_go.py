"""
Doc: https://www.pantsbuild.org/docs/plugins-codegen
The protobuf build for Go files are under development by pants.
The following is a simplified version of the go code generation for protobuf files.
Please follow https://github.com/pantsbuild/pants/tree/main/src/python/pants/backend/codegen/protobuf/go
"""
import os

from pants.engine.console import Console
from pants.engine.fs import Workspace, PathGlobs
from pants.engine.goal import GoalSubsystem, Goal
from pants.engine.internals.native_engine import Digest, Snapshot, MergeDigests
from pants.engine.internals.selectors import Get
from pants.engine.process import ProcessResult, Process, BinaryPaths, BinaryPathRequest
from pants.engine.rules import goal_rule, collect_rules


class ProtoCompileGoSubsystem(GoalSubsystem):
    name = "proto-compile-go"
    help = "A goal to compile proto files into go files."


class ProtoCompileGo(Goal):
    subsystem_cls = ProtoCompileGoSubsystem


@goal_rule
async def proto_compile_go(console: Console, workspace: Workspace) -> ProtoCompileGo:
    proto_digest = await Get(Digest, PathGlobs(["rpc/**/*.proto"]))
    proto_snapshot = await Get(Snapshot, Digest, proto_digest)

    bin_digest = await Get(Digest, PathGlobs(["bin/**"]))
    bin_snapshot = await Get(Snapshot, Digest, bin_digest)

    output_files = []
    for path in proto_snapshot.files:
        output_files.append(os.path.splitext(path)[0]+'.twirp.go')
        output_files.append(os.path.splitext(path)[0]+'.pb.go')

    protoc_paths = await Get(
        BinaryPaths,
        BinaryPathRequest(
            binary_name="protoc",
            search_path=["/usr/local/bin", "/opt/homebrew/bin", "/usr/bin", "/bin", "<PATH>"],
        )
    )
    protoc_bin = protoc_paths.first_path
    if protoc_bin is None:
        raise OSError("Could not find binary 'protoc'.")

    result = await Get(
        ProcessResult,
        Process(
            argv=[protoc_bin.path, "--twirp_out=.", "--go_out=.", *proto_snapshot.files],
            description="Generate go files for all proto files",
            input_digest=await Get(Digest, MergeDigests([proto_snapshot.digest, bin_snapshot.digest])),
            output_files=output_files,
            env={"PATH": "bin"},
        )
    )

    if result.stderr.decode() is not None and result.stderr.decode() != "":
        console.print_stdout(result.stderr.decode())
        return ProtoCompileGo(exit_code=1)

    workspace.write_digest(result.output_digest)

    return ProtoCompileGo(exit_code=0)


def rules():
    return collect_rules()
