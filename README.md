# bridge-platform

## Backend

```
cd backend
python -m venv venv
source venv/bin/activate
pip install pipenv
pipenv install --dev
pipenv run python manage.py migrate
pipenv run python manage.py runserver
```

## Frontend

```
cd frontend
npm install
npm run dev
```