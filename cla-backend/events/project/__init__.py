import os
from flask import Flask



def create_app(script_info=None):

    app = Flask(__name__)
    app_settings = os.getenv('APP_SETTINGS')
    app.config.from_object(app_settings)

    #register blueprints
    from project.api.events import events_blueprint
    app.register_blueprint(events_blueprint)

    def ctx():
        return {'app': app}

    return app 

