FROM postgres:latest
ENV POSTGRES_USER kanban_user 
ENV POSTGRES_PASSWORD kanban_password 
ENV POSTGRES_DB kanban_db
RUN chmod 644 ./migrations/*.sql
RUN apt-get update && apt-get install -y postgresql-uuid-ossp