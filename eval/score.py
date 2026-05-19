#!/usr/bin/env python3
"""Eval script for rta — Go-centric project eval.

Dimensions: tests (go test), coverage (go test -cover), lint (go vet), observability (source scan).
"""

import json
import re
import subprocess
import sys


def eval_tests() -> dict:
    """Run test suite: go test ./..."""
    try:
        result = subprocess.run(
            ['go', 'test', './...'],
            capture_output=True,
            text=True,
            timeout=120,
        )
        passed = result.returncode == 0
        if passed:
            score = 1.0
        else:
            error_lines = [ln for ln in (result.stdout + result.stderr).splitlines() if ln.strip()]
            score = max(0.0, 1.0 - len(error_lines) * 0.05) if error_lines else 0.0
        return {
            "name": "tests",
            "score": score,
            "weight": 0.30,
            "passed": passed,
            "details": (result.stdout or result.stderr).strip()[-500:],
        }
    except subprocess.TimeoutExpired:
        return {"name": "tests", "score": 0.0, "weight": 0.30, "passed": False, "details": "Timed out after 120s"}
    except FileNotFoundError:
        return {"name": "tests", "score": 0.0, "weight": 0.30, "passed": False, "details": "Go toolchain not available"}


def eval_coverage() -> dict:
    """Measure test coverage: go test -cover ./..."""
    try:
        result = subprocess.run(
            ['go', 'test', '-cover', './...'],
            capture_output=True,
            text=True,
            timeout=120,
        )
        coverages = []
        for line in result.stdout.splitlines():
            m = re.search(r'coverage:\s+([\d.]+)%', line)
            if m:
                coverages.append(float(m.group(1)))
            elif 'no test files' in line:
                coverages.append(0.0)

        if not coverages:
            return {"name": "coverage", "score": 0.0, "weight": 0.40, "passed": False, "details": "No coverage data"}

        avg = sum(coverages) / len(coverages)
        score = avg / 100.0
        passed = avg >= 95.0
        pkg_details = ", ".join(f"{c:.0f}%" for c in coverages)
        return {
            "name": "coverage",
            "score": round(score, 4),
            "weight": 0.40,
            "passed": passed,
            "details": f"avg={avg:.1f}% packages=[{pkg_details}]",
        }
    except subprocess.TimeoutExpired:
        return {"name": "coverage", "score": 0.0, "weight": 0.40, "passed": False, "details": "Timed out after 120s"}
    except FileNotFoundError:
        return {"name": "coverage", "score": 0.0, "weight": 0.40, "passed": False, "details": "Go toolchain not available"}


def eval_lint() -> dict:
    """Run linter: go vet ./..."""
    try:
        result = subprocess.run(
            ['go', 'vet', './...'],
            capture_output=True,
            text=True,
            timeout=120,
        )
        passed = result.returncode == 0
        if passed:
            score = 1.0
        else:
            error_lines = [ln for ln in (result.stdout + result.stderr).splitlines() if ln.strip()]
            score = max(0.0, 1.0 - len(error_lines) * 0.05) if error_lines else 0.0
        return {
            "name": "lint",
            "score": score,
            "weight": 0.15,
            "passed": passed,
            "details": (result.stdout or result.stderr).strip()[-500:],
        }
    except subprocess.TimeoutExpired:
        return {"name": "lint", "score": 0.0, "weight": 0.15, "passed": False, "details": "Timed out after 120s"}
    except FileNotFoundError:
        return {"name": "lint", "score": 0.0, "weight": 0.15, "passed": False, "details": "Go toolchain not available"}


def eval_observability() -> dict:
    """Analyze observability coverage: logging, structured logging, request tracing in Go code."""
    from pathlib import Path

    skip = {"vendor", ".git", ".factory", "eval", "testdata"}
    log_pats = [r"\blog\.\w+\(", r"\bfmt\.Printf?\(", r"\bfmt\.Fprintf?\(", r"\bslog\.\w+\("]
    struct_pats = [r"\bslog\b", r"\bzerolog\b", r"\bzap\b", r"\blogrus\b"]
    trace_pats = [r"\bopentelemetry\b", r"\btrace\b", r"\bspan\b"]

    sources = [f for f in Path(".").rglob("*.go")
               if not any(p in f.parts for p in skip)
               and not f.name.endswith("_test.go")]
    total_fn = logged_fn = total_log = 0
    has_struct = has_trace = False
    func_pat = re.compile(r"^func\s+(?:\([^)]+\)\s+)?\w+\(")

    for src in sources:
        try:
            code = src.read_text(errors="replace")
        except OSError:
            continue
        lines = code.splitlines()
        in_func = False
        func_has_log = False
        for line in lines:
            if func_pat.match(line):
                if in_func and func_has_log:
                    logged_fn += 1
                in_func = True
                func_has_log = False
                total_fn += 1
            elif in_func:
                for pat in log_pats:
                    if re.search(pat, line):
                        func_has_log = True
                        total_log += 1
                        break
        if in_func and func_has_log:
            logged_fn += 1
        for pat in struct_pats:
            if re.search(pat, code):
                has_struct = True
        for pat in trace_pats:
            if re.search(pat, code, re.IGNORECASE):
                has_trace = True

    if total_fn == 0:
        return {"name": "observability", "score": 0.0, "weight": 0.15, "passed": True, "details": "No functions found"}

    cov = logged_fn / total_fn
    density = min(1.0, total_log / max(total_fn, 1))
    score = 0.40 * cov + 0.25 * float(has_struct) + 0.20 * float(has_trace) + 0.15 * density
    details = (f"coverage={cov:.0%} ({logged_fn}/{total_fn}), "
               f"structured={'yes' if has_struct else 'no'}, "
               f"tracing={'yes' if has_trace else 'no'}, "
               f"density={density:.0%}")

    return {"name": "observability", "score": round(score, 3), "weight": 0.15, "passed": score >= 0.3, "details": details}


EVALS = [eval_tests, eval_coverage, eval_lint, eval_observability]


def main() -> None:
    results = [fn() for fn in EVALS]
    output = {"results": results}
    json.dump(output, sys.stdout, indent=2)
    print()


if __name__ == "__main__":
    main()
