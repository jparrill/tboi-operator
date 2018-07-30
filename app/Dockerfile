FROM python:2.7-alpine
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt && \
    mkdir ./templates
COPY tboi_items_server.py tboi_items_ref.py ./
COPY templates/index.html.j2 ./templates
CMD [ "python", "./tboi_items_server.py" ]
