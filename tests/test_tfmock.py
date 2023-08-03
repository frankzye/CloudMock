import os

import pytest
import tftest
from tfmock import TFMock


@pytest.fixture(scope="class")
def plan():
    tf_mock = TFMock()
    try:
        tf = tftest.TerraformTest('azure', basedir=os.path.dirname(__file__))
        tf_mock.setup(tf, cleanup_on_exit=False)
        yield tf.plan(output=True)
    finally:
        tf_mock.clean_up()


class TestCases:
    def test_resources(self, plan: tftest.TerraformPlanOutput):
        rg = plan.resources["azurerm_resource_group.rg"]
        assert rg['values']['name'] == 'xx'
