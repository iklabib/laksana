from types import ModuleType, TracebackType
from typing import Iterable, TextIO
import unittest
from main import hello

class Test(unittest.TestCase):
    def test_ini(self):
        ret = hello()
        exp = 'Hello, World!'
        self.assertEqual(ret, exp, msg=f'Kembalian {ret}, kembalian yang diharapkan adalah "{exp}"')
    
    def get_test_name(self):
        return "INI"

class TestResult(unittest.TestResult):
    def __init__(self, stream: TextIO | None = None, descriptions: bool | None = None, verbosity: int | None = None) -> None:
        super().__init__(stream, descriptions, verbosity)
        self.success = []
    
    def addSuccess(self, test: Test) -> None:
        self.success.append(test.get_test_name())
        return super().addSuccess(test)

class TestRunner(unittest.TextTestRunner):
    def _makeResult(self) -> unittest.TestResult:
        return TestResult(self.stream, self.descriptions, self.verbosity)
    
    def run(self, test: unittest.TestSuite | unittest.TestCase) -> unittest.TestResult:
        return super().run(test)

if __name__ == '__main__':
    runner = TestRunner()
    ut = unittest.main(testRunner=runner, verbosity=0, exit=False)