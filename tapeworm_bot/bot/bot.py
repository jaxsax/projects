import lib

if __name__ == "__main__":
    app = lib.create_app()
    app.run('127.0.0.1', port=8080, debug=True)
