# bridge-platform

## Backend

```
cd backend
python3 -m venv .venv
source .venv/bin/activate # For windows use: source venv\Scripts\activate
pip install pipenv
pipenv install --dev
pipenv run python manage.py migrate
pipenv run python manage.py runserver
deactivate # To kill venv
```

## Frontend

```
cd frontend
npm install
npm run dev
```