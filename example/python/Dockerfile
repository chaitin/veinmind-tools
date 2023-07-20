FROM veinmind/python3.6:1.9.42-stretch
WORKDIR /tool
ADD . .
RUN pip install -r requirements.txt
RUN echo "#!/bin/bash\n\npython scan.py \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh
ENTRYPOINT ["/tool/entrypoint.sh"]

