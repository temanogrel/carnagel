from setuptools import setup, find_packages

setup(
    name='common',
    version=4.0,
    author='Palmer Raynard',
    packages=find_packages(),
    install_requires=[
        'cffi',
        'python-consul',
        'grpcio',
        'grpcio-tools',
        'requests>=2.5.0',
        'nose>=1.3.7',
        'iso8601>=0.1'
    ],
    dependency_links=[
        'git+ssh://git@git.misc.vee.bz/carnagel/minerva-bindings.git@0.0.2#egg=minerva'
    ]
)
