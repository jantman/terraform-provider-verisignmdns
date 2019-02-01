from flask import Flask, g, request
from flask_restful import Resource, Api, abort, reqparse
from datetime import datetime, timedelta

app = Flask(__name__)
api = Api(app)


ACCOUNT_ID = 9999999
ZONE_NAME = 'example.com'
ZONE_ID = 123456
AUTH_TOKEN = '324f3919c575bae096c4df3e638a83d6'

parser = reqparse.RequestParser()
parser.add_argument('owner', type=str)
parser.add_argument('view_name', type=str)
parser.add_argument('type', type=str)
parser.add_argument('rdata', type=str)
parser.add_argument('ttl', type=int)
parser.add_argument('comments', type=str)

dts = (datetime.now() - timedelta(hours=27)).strftime(
    '%Y-%m-%dT%H:%M:%S.000%z'
)
RECORDS = {
    1234: {
        "resource_record_id": "1234",
        'modifiable': True,
        "owner": "www.example.com",
        "type": "A",
        "rdata": "192.0.2.29",
        "ttl": "86400",
        "create_timestamp": dts,
        "update_timestamp": dts,
        "created_by": "someuser",
        "updated_by": "someuser"
    }
}

def abort_if_bad_account_or_zone(accountId, zoneName):
    if accountId != ACCOUNT_ID:
        app.logger.info('401 - wrong account ID')
        abort(
            401,
            error_code="401",
            error_messages=["You are not authorized for this call."]
        )
    if zoneName != ZONE_NAME:
        app.logger.info('422 - wrong zone name')
        abort(
            422,
            error_code="ERROR_RULE_VIOLATION",
            error_messages=["Zone not found"]
        )


def abort_if_bad_auth():
    auth = request.headers.get('Authorization', '')
    if auth != 'Token %s' % AUTH_TOKEN:
        app.logger.info('401 - wrong token')
        abort(
            401,
            error_code="401",
            error_messages=["Your credentials are not valid."]
        )


def abort_if_bad_headers():
    if request.headers.get('Content-Type', '') != 'application/json':
        app.logger.info('400 - wrong Content-Type')
        abort(400, message='bad headers')
    if request.headers.get('Accept', '') != 'application/json':
        app.logger.info('400 - wrong Accept')
        abort(400, message='bad headers')


class ResourceRecord(Resource):

    def get(self, accountId, zoneName, resourceRecordId):
        abort_if_bad_auth()
        abort_if_bad_headers()
        app.logger.info(vars(request))
        return {'hello': 'world'}

    def put(self, accountId, zoneName, resourceRecordId):
        abort_if_bad_auth()
        abort_if_bad_headers()
        RECORDS[resourceRecordId] = request.form['data']
        return {}, 201, {'X-Header': 'value'}


class ResourceRecordList(Resource):

    def get(self, accountId, zoneName):
        abort_if_bad_auth()
        abort_if_bad_headers()
        abort_if_bad_account_or_zone(accountId, zoneName)
        return {
            'total_count': len(RECORDS),
            'zone_name': ZONE_NAME,
            'zone_id': ZONE_ID,
            'view_name': 'DEFAULT',
            'resource_records': [list(RECORDS.values())]
        }, 200

    def post(self, accountId, zoneName):
        abort_if_bad_auth()
        abort_if_bad_headers()
        abort_if_bad_account_or_zone(accountId, zoneName)
        args = parser.parse_args(strict=True)
        app.logger.info('POST: %s' % args)
        if args['owner'] == 'bad.example.com':
            app.logger.info('POST - return 400 for bad.example.com')
            abort(
                400,
                error_code='MALFORMED', error_messages=['foo']
            )
        if (
            'owner' not in args or
            'rdata' not in args or
            'type' not in args
        ):
            app.logger.info('POST - 422 - missing required fields')
            abort(
                422,
                error_code='BAD_FIELDS', error_messages=['foo']
            )
        k = max(RECORDS.keys()) + 1
        dts = datetime.now().strftime('%Y-%m-%dT%H:%M:%S.000%z')
        RECORDS[k] = {
            "resource_record_id": '%d' % k,
            'modifiable': True,
            "owner": args['owner'],
            "type": args['type'],
            "rdata": args['rdata'],
            "ttl": args.get('ttl', 86400),
            "create_timestamp": dts,
            "update_timestamp": dts,
            "created_by": "someuser",
            "updated_by": "someuser"
        }
        loc = 'http://%s/api/v1/accounts/%s/zones/%s/rr/%d' % (
            request.environ['HTTP_HOST'], ACCOUNT_ID, ZONE_NAME, k
        )
        return '', 201, {'Location': loc}


api.add_resource(
    ResourceRecord,
    '/api/v1/accounts/<int:accountId>/zones/<string:zoneName>/rr/<int:resourceRecordId>'
)
api.add_resource(
    ResourceRecordList,
    '/api/v1/accounts/<int:accountId>/zones/<string:zoneName>/rr'
)


if __name__ == '__main__':
    app.run(debug=True, threaded=False)
