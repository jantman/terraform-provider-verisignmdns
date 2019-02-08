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
        """
        GET /api/v1/accounts/{accountId}/zones/{zoneName}/rr/{resourceRecordId}
        Accept: application/json
        Content-Type: application/json

        200
        {
          "resource_record_id": "1234",
          "owner": "www.example.com",
          "type": "A",
          "rdata": "192.0.2.29",
          "ttl": "86400",
          "create_timestamp": "2011-01-21T16:35:35.000+00:00",
          "update_timestamp": "2011-01-21T16:35:35.000+00:00",
          "created_by": "someuser",
          "updated_by": "someotheruser"
        }

        400 - The request was malformed. Please check the format of the request.

        422 - The request was well-formed, but one or more validations failed.
        {
          "error_code": "ERROR_OPERATION_FAILURE",
          "error_messages": [
            "Domain already exists. Please verify your domain name."
          ]
        }
        """
        abort_if_bad_auth()
        abort_if_bad_headers()
        abort_if_bad_account_or_zone(accountId, zoneName)
        if resourceRecordId not in RECORDS:
            app.logger.info('DELETE - 422 - bad ID')
            abort(
                422,
                error_code='BAD_ID', error_messages=['record ID invalid']
            )
        return RECORDS[resourceRecordId], 200

    def put(self, accountId, zoneName, resourceRecordId):
        """
        PUT /api/v1/accounts/{accountId}/zones/{zoneName}/rr/{resourceRecordId}
        Accept: application/json
        Content-Type: application/json

        Request:
        {
          "rdata": "192.0.2.30",
          "comments": "Updating a resource record."
        }

        200 Resource record successfully updated.
        Headers:
        Location - URL of created resource record with its resource record Id

        400 - The request was malformed. Please check the format of the request.

        422 - The request was well-formed, but one or more validations failed.
        {
          "error_code": "ERROR_OPERATION_FAILURE",
          "error_messages": [
            "Domain already exists. Please verify your domain name."
          ]
        }
        """
        abort_if_bad_auth()
        abort_if_bad_headers()
        abort_if_bad_account_or_zone(accountId, zoneName)
        args = parser.parse_args(strict=True)
        app.logger.info('PUT: %s' % args)
        if not set(args.keys).issubset(set(['rdata', 'comments'])):
            app.logger.info('PUT - 422 - bad fields')
            abort(
                422,
                error_code='BAD_FIELDS', error_messages=['foo']
            )
        if resourceRecordId not in RECORDS:
            app.logger.info('PUT - 422 - bad ID')
            abort(
                422,
                error_code='BAD_ID', error_messages=['record ID invalid']
            )
        if 'rdata' in args:
            RECORDS[resourceRecordId]['rdata'] = args['rdata']
        if 'comments' in args:
            RECORDS[resourceRecordId]['comments'] = args['comments']
        loc = 'http://%s/api/v1/accounts/%s/zones/%s/rr/%d' % (
            request.environ['HTTP_HOST'], ACCOUNT_ID, ZONE_NAME, k
        )
        return '', 200, {'Location': loc}

    def delete(self, accountId, zoneName, resourceRecordId):
        """
        DELETE /api/v1/accounts/{accountId}/zones/{zoneName}/rr/{resourceRecordId}
        Accept: application/json
        Content-Type: application/json

        Request:
        {
          "comments": "Deleting example.com zone"
        }

        204 - Successfully deleted

        400 - The request was malformed. Please check the format of the request.

        422 - The request was well-formed, but one or more validations failed.
        {
          "error_code": "ERROR_OPERATION_FAILURE",
          "error_messages": [
            "Domain already exists. Please verify your domain name."
          ]
        }
        """
        abort_if_bad_auth()
        abort_if_bad_headers()
        abort_if_bad_account_or_zone(accountId, zoneName)
        args = parser.parse_args(strict=True)
        app.logger.info('DELETE: %s' % args)
        if len(args.keys) > 0 and list(args.keys) != ['comments']:
            app.logger.info('DELETE - 422 - bad fields')
            abort(
                422,
                error_code='BAD_FIELDS', error_messages=['foo']
            )
        if resourceRecordId not in RECORDS:
            app.logger.info('DELETE - 422 - bad ID')
            abort(
                422,
                error_code='BAD_ID', error_messages=['record ID invalid']
            )
        del RECORDS[resourceRecordId]
        return '', 204


class ResourceRecordList(Resource):

    def get(self, accountId, zoneName):
        """
        GET /api/v1/accounts/{accountId}/zones/{zoneName}/rr
        Accept: application/json
        Content-Type: application/json

        API is paginated, but this mock doesn't expose that, as the provider
        doesn't use this endpoint

        200 OK:
        {
          "total_count": 2,
          "zone_name": "example.com",
          "zone_id": 12345,
          "view_name": "ASIA",
          "resource_records":[
                {
                  "resource_record_id": "1234",
                  "modifiable": true,
                  "owner": "www.example.com",
                  "type": "A",
                  "rdata": "192.0.2.25",
                  "ttl": "86400",
                  "create_timestamp": "2011-01-21T16:35:35.000+00:00",
                  "update_timestamp": "2011-01-21T16:35:35.000+00:00",
                  "created_by": "someuser",
                  "updated_by": "someotheruser"
                },
                {
                  "resource_record_id": "1235",
                  "modifiable": true,
                  "owner": "www1.example.com",
                  "type": "A",
                  "rdata": "192.0.2.26",
                  "ttl": "86400",
                  "create_timestamp": "2011-01-21T16:35:35.000+00:00",
                  "update_timestamp": "2011-01-21T16:35:35.000+00:00",
                  "created_by": "someuser",
                  "updated_by": "someotheruser"
                }
            ]
        }

        400 - The request was malformed. Please check the format of the request.

        422 - The request was well-formed, but one or more validations failed.
        {
          "error_code": "ERROR_OPERATION_FAILURE",
          "error_messages": [
            "Domain already exists. Please verify your domain name."
          ]
        }
        """
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
        """
        POST /api/v1/accounts/{accountId}/zones/{zoneName}/rr
        Accept: application/json
        Content-Type: application/json

        Request:
        {
          "owner": "www.example.com",
          "type": "A",
          "rdata": "192.0.2.27",
          "comments": "Adding a resource record"
        }

        201 Resource record successfully created.
        Headers:
        Location - URL of created resource record with its resource record Id

        400 - The request was malformed. Please check the format of the request.

        422 - The request was well-formed, but one or more validations failed.
        {
          "error_code": "ERROR_OPERATION_FAILURE",
          "error_messages": [
            "Domain already exists. Please verify your domain name."
          ]
        }
        """
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
