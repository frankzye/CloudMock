import asyncio
import os
import stat
import sys
from importlib.metadata import version

import tftest


class TFMock:
    def __init__(self):
        self.proc = None
        self.mock_root = None

    def setup(self, tf: tftest.TerraformTest, **kwargs):
        top_root = os.path.dirname(__file__)
        mock_root = os.path.join(top_root, f'{__name__}-{version(__name__)}.dist-info', 'bin')
        self.mock_root = mock_root

        os.environ["GIT_SSL_NO_VERIFY"] = "true"
        os.environ["cloud_mock_root"] = mock_root

        os.chmod(os.path.join(mock_root, sys.platform, "az"), stat.S_IEXEC)
        os.chmod(os.path.join(mock_root, sys.platform, "CloudMock"), stat.S_IEXEC)

        tf.env.update(os.environ)

        asyncio.run(self.mock_server())

        tf.setup(**kwargs)

        os.environ["EQUESTS_CA_BUNDLE"] = os.path.join(mock_root, "public.pem")
        os.environ["SSL_CERT_FILE"] = os.path.join(mock_root, "public.pem")
        os.environ["HTTPS_PROXY"] = "http://127.0.0.1:9999"
        if sys.platform == "win32":
            os.environ["PATH"] = os.path.join(self.mock_root, sys.platform) + ";" + os.environ["PATH"]
        else:
            os.environ["PATH"] = os.path.join(self.mock_root, sys.platform) + ":" + os.environ["PATH"]

        tf.env.update(os.environ)

    async def mock_server(self):
        if sys.platform == "win32":
            cmd = os.path.join(self.mock_root, sys.platform, "CloudMock.exe")
        else:
            cmd = os.path.join(self.mock_root, sys.platform, "CloudMock")
        self.proc = await asyncio.create_subprocess_exec(cmd)

    def clean_up(self):
        if self.proc is not None:
            self.proc.terminate()
