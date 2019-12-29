import sys
import pytest

import lib as l
from lib import name_builder


def test_works():
    assert 1 == 1


def test_name_builder():
    got = name_builder("a", "b")
    want = "a b"
    assert want == got


def test_name_builder_import():
    got = l.name_builder("a", "b")
    want = "a b"
    assert want == got


if __name__ == "__main__":
    sys.exit(pytest.main([__file__]))
