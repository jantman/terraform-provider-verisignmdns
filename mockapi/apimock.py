from flask import Flask, g
from flask_restful import Resource, Api

app = Flask(__name__)
api = Api(app)


ACCOUNT_ID = 9999999
ZONE_NAME = 'example.com'
ZONE_ID = 123456


def get_data():
    if 'data' in g:
        return g.data
    g.data = {
        1234: {
            "resource_record_id": "1234",
            'modifiable': True,
            "owner": "www.example.com",
            "type": "A",
            "rdata": "192.0.2.29",
            "ttl": "86400",
            "create_timestamp": "2011-01-21T16:35:35.000+00:00",
            "update_timestamp": "2011-01-21T16:35:35.000+00:00",
            "created_by": "someuser",
            "updated_by": "someotheruser"
        }
    }


def abort_if_bad_account_or_zone(accountId, zoneName):
    if accountId != ACCOUNT_ID:
        abort(
            401,
            {
                "error_code": "401",
                "error_messages": ["You are not authorized for this call."]
            },
            {
                'X-Frame-Options': 'DENY',
                'Strict-Transport-Security': 'max-age=15768000',
                'Expires': '-1',
                'Cache-Control': 'private, no-store, no-cache, must-revalidate, max-age=0, proxy-revalidate, s-maxage=0',
                'Pragma': 'no-cache',
                'Cache-Control': 'post-check=0, pre-check=0',
                'WWW-Authenticate': 'Basic realm="MDNS"',
                'Vary': 'Accept-Encoding',
                'Content-Type': 'application/octet-stream'
            }
        )
    if zoneName != ZONE_NAME:
        abort(
            422,
            {
                "error_code":"ERROR_RULE_VIOLATION",
                "error_messages":["Zone not found"]
            },
            {
                'X-Frame-Options': 'DENY',
                'Strict-Transport-Security': 'max-age=15768000',
                'Expires': '-1',
                'Cache-Control': 'private, no-store, no-cache, must-revalidate, max-age=0, proxy-revalidate, s-maxage=0',
                'Pragma': 'no-cache',
                'Cache-Control': 'post-check=0, pre-check=0',
                'Vary': 'Accept-Encoding',
                'Content-Type': 'application/json'
            }
        )


class ResourceRecord(Resource):

    def get(self, accountId, zoneName, resourceRecordId):
        return {'hello': 'world'}

    def put(self, accountId, zoneName, resourceRecordId):
        g.data[resourceRecordId] = request.form['data']
        return {}, 201, {'X-Header': 'value'}


class ResourceRecordList(Resource):

    def get(self, accountId, zoneName):
        abort_if_bad_account_or_zone(accountId, zoneName)
        get_data()
        return {
            'total_count': len(g.data),
            'zone_name': ZONE_NAME,
            'zone_id': ZONE_ID,
            'view_name': 'DEFAULT',
            'resource_records': [list(g.data.values())]
        }, 200

    def post(self, accountId, zoneName):
        abort_if_bad_account_or_zone(accountId, zoneName)
        get_data()


api.add_resource(ResourceRecord, '/api/v1/accounts/<int:accountId>/zones/<string:zoneName>/rr/<int:resourceRecordId>')
api.add_resource(ResourceRecordList, '/api/v1/accounts/<int:accountId>/zones/<string:zoneName>/rr')


if __name__ == '__main__':
    app.run(debug=True, threaded=False)
