from flask import Flask

def create_app(env='dev'):
    app = Flask(
        __name__
    )

    @app.route('/')
    def root():
        return {'a': 1}

    return app