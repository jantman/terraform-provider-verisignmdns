# Mock Verisign MDNS API

Verisign MDNS doesn't provide an API sandbox or anywhere for testing. The code in this
directory uses Python's [flask-restful](https://flask-restful.readthedocs.io/en/latest/index.html)
to create a really simple mock of the Verisign MDNS API.

## Usage

```bash
python3 -mvenv .venv
source .venv/bin/activate
pip install flask-restful
./apimock.py
```
