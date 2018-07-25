find . -type d -name "__pycache__" -exec rm -rf {} +
find . -type f -name "*.py[co]" -delete
find . -type f -name "*.log" -delete
