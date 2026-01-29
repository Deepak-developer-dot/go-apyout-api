
- name: Build and Push Docker Image
  uses: docker/build-push-action@v5
  with:
    context: ./backend
    push: true
    tags: 9534/payout-api:latest
