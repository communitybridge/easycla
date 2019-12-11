import json
import uuid


def test_add_event(test_app, events_table):
    """
    Test creating events instance in dynamodb 
    """
    event_id = str(uuid.uuid4())
    client = test_app.test_client()
    response = client.post(
        "/events",
        data=json.dumps({"event_id": event_id}),
        content_type="application/json",
    )
    data = json.loads(response.data.decode())
    assert response.status_code == 200
    assert "success" in data["status"]



