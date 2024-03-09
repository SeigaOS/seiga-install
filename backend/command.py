import subprocess


def run(cmd: str | list[str], can_err=True) -> str:
    if isinstance(cmd, str):
        res = subprocess.run(cmd, shell=True, check=True, capture_output=True,)
    else:
        res = subprocess.run(cmd, shell=False, check=True, capture_output=True)

    out = res.stdout.decode("utf-8")
    if can_err and res.returncode != 0:
        raise RuntimeError(
            f"Command: {cmd} failed with return code {res.returncode}\n{out}"
        )
    return out
