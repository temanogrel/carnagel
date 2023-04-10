<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording;

class DeathFilePermissions
{
    // CRUD death file permissions
    const VIEW_DEATH_FILE = 'aphrodite:recording-death-files:view';
    const LIST_DEATH_FILES = 'aphrodite:recording-death-files:list';
    const UPLOAD_DEATH_FILE = 'aphrodite:recording-death-files:upload';
    const UPDATE_DEATH_FILE = 'aphrodite:recording-death-files:update';
    const DELETE_DEATH_FILE = 'aphrodite:recording-death-files:delete';

    // CRUD death file url permissions
    const ADD_ENTRY = 'aphrodite:recording-death-file:url:add';
    const VIEW_ENTRY = 'aphrodite:recording-death-file:url:view';
    const LIST_ENTRIES = 'aphrodite:recording-death-file:url:list';
    const REMOVE_ENTRY = 'aphrodite:recording-death-file:url:remove';
    const UPDATE_ENTRY = 'aphrodite:recording-death-file:url:update';
}
