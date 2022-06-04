# Pending

- [ ] Delete all the old code (things that aren't launched from simplifiedbot)
- [ ] Cleanup all the bazel stuff it's not being used anymore
- [ ] Domain blacklist
- [ ] Whitelist messages only from certain chat areas
- [ ] Some ways to label links around topics??
  - Only web, or from the chat?

# Mess around

- [ ] Migrate the storage to use foundationdb
- [ ] Implement some form of search with foundationdb
  - Let's say some link got deleted, how to remove it from the index?
  - given this KV structure; (index_space, word, docid) -> 1
  - Or you store the deleted docids somewhere like (deleted_space, docid) -> 1
  - And when the application reads the responses, it filters out the docid. But this still takes up space in the index, still need some way to remove it from the index periodically (investigate how elasticsearch -> lucene does this with the tombstone approach)