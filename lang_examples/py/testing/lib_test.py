import unittest

import lib as l
from lib import name_builder


class TypicalUnitTest(unittest.TestCase):
    def test_works(self):
        self.assertEqual(1, 1)

    def test_name_builder(self):
        got = name_builder("a", "b")
        want = "a b"
        self.assertEqual(want, got)

    def test_name_builder_import(self):
        got = l.name_builder("a", "b")
        want = "a b"
        self.assertEqual(want, got)
