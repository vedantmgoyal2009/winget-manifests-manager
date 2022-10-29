from http.server import BaseHTTPRequestHandler
from sqlite3 import connect
from json import dumps
from zipfile import ZipFile
from urllib.request import urlretrieve

# Download the source.msix file from the WinGet CDN (Content Delivery Network)
urlretrieve("https://cdn.winget.microsoft.com/cache/source.msix", "/tmp/source.msix")

# Extract the index.db file from the source.msix file
ZipFile("/tmp/source.msix").extract("Public/index.db", "/tmp/")

# Connect to the database
db = connect("/tmp/Public/index.db")
cursor = db.cursor()

# Get the manifests table
cursor.execute('SELECT id, version FROM manifest')
manifests = cursor.fetchall()

# Initialize the result dictionary
result = {}

# Iterate over each row in the manifests table
for row in manifests:
    id, version = row

    # Get the id value from the ids table
    cursor.execute(f'SELECT id FROM ids WHERE rowid = {id}')
    id_value = cursor.fetchone()[0]

    # Get the version value from the versions table
    cursor.execute(f'SELECT version FROM versions WHERE rowid = {version}')
    version_value = cursor.fetchone()[0]

    # Add the id and version to the result dictionary
    if id_value not in result:
        result[id_value] = []
    result[id_value].append(version_value)

# Close the database connection
db.close()

print(dumps(result).encode())

class handler(BaseHTTPRequestHandler):

  def do_GET(self):
    self.send_response(200)
    self.send_header('Content-type', 'text/json')
    self.end_headers()
    self.wfile.write(dumps(result).encode())
    return
