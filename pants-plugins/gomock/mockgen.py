"""
Doc: https://www.pantsbuild.org/docs/plugins-codegen
The protobuf build for Go files are under development by pants.
The following is a simplified version of the go code generation for protobuf files.
Please follow https://github.com/pantsbuild/pants/tree/main/src/python/pants/backend/codegen/protobuf/go
"""
import os

from pants.engine.console import Console
from pants.engine.fs import Workspace, PathGlobs, CreateDigest, FileContent
from pants.engine.goal import GoalSubsystem, Goal
from pants.engine.internals.native_engine import Digest, Snapshot, MergeDigests
from pants.engine.internals.selectors import Get, MultiGet
from pants.engine.process import ProcessResult, Process
from pants.engine.rules import goal_rule, collect_rules


class GoMockSubsystem(GoalSubsystem):
    name = "gomock"
    help = "A goal to generate mock files for go interface"


class GoMock(Goal):
    subsystem_cls = GoMockSubsystem


@goal_rule
async def mockgen(console: Console, workspace: Workspace) -> GoMock:
    go_file_digest = await Get(Digest, PathGlobs(["**/*.go", "!**/*_test.go"]))
    go_file_snapshot = await Get(Snapshot, Digest, go_file_digest)

    output_file_attrs = []
    for path in go_file_snapshot.files:
        output_file_attrs.append((
            path,
            os.path.join('mocks', path),
            os.path.basename(os.path.dirname(path))
        ))

    bin_digest = await Get(Digest, PathGlobs(["bin/**"]))
    bin_snapshot = await Get(Snapshot, Digest, bin_digest)

    gomod_digest = await Get(Digest, PathGlobs(["go.mod"]))
    gomod_snapshot = await Get(Snapshot, Digest, gomod_digest)

    env = {"PATH": "bin", "GOBIN": "bin", "GOMOD": "go.mod"}

    tasks = []
    for i in range(len(go_file_snapshot.files)):
        tasks.append(Get(ProcessResult,
                         Process(
                             argv=[
                                 "mockgen",
                                 "-source", output_file_attrs[i][0],
                                 "-package", output_file_attrs[i][2]
                             ],
                             description="Generate go mocks",
                             input_digest=await Get(Digest, MergeDigests([
                                 go_file_snapshot.digest, bin_snapshot.digest, gomod_snapshot.digest])),
                             env=env,
                         )))

    result = await MultiGet(*tasks)

    has_error = False
    outputs = []
    for i, each in enumerate(result):
        if each.stderr.decode() is not None and each.stderr.decode() != "":
            has_error = True
            console.print_stdout("failed to generate mock file for {}, error: {}".format(
                output_file_attrs[i][0], each.stderr.decode()))
        else:
            gen = str(each.stdout.decode()).strip().split('\n')[-1]
            if not gen.endswith("package {}".format(output_file_attrs[i][2])):
                outputs.append(FileContent(output_file_attrs[i][1], str.encode(each.stdout.decode())))

    workspace.write_digest(await Get(Digest, CreateDigest(outputs)))

    return GoMock(exit_code=0 if not has_error else 1)


def rules():
    return collect_rules()
