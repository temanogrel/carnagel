def deep_update(orig_dict: dict, new_dict: dict):
    for key, val in new_dict.items():

        if isinstance(val, dict):
            tmp = deep_update(orig_dict.get(key, {}), val)
            orig_dict[key] = tmp
        elif isinstance(val, list):
            orig_dict[key] = (orig_dict[key] + val)
        else:
            orig_dict[key] = new_dict[key]

    return orig_dict


def get_recordings_matching(max_results=None, **criteria):
    """
    Retrieve all recordings, or until we hit the max results of recordings matching the given criteria

    :param max_results:
    :param criteria:

    :return:
    """

    recording_api = recording_api_factory('central')

    limit = 500
    offset = 0
    yielded = 0

    while True:

        result = recording_api.get_all(limit=limit, offset=offset, **criteria)

        for recording in result.items:
            yield recording

            yielded += 1

            if max_results is not None and max_results <= yielded:
                break

        # Short circuit
        if max_results is not None and max_results <= yielded:
            break

        # No more recordings left
        if offset > result.total:
            break

        offset += limit
