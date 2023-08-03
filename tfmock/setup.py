import os.path
import pathlib

from setuptools import setup

root_folder = os.path.dirname(__file__)
long_description = pathlib.Path(os.path.join(root_folder, "readme.md")).read_text()

setup(
    name='tfmock',
    version='0.0.2',
    install_requires=['tftest'],
    include_package_data=True,
    long_description=long_description,
    long_description_content_type='text/markdown'
)
