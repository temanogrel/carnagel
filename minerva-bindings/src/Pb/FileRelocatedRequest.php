<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: file.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Protobuf type <code>pb.FileRelocatedRequest</code>
 */
class FileRelocatedRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * <code>string uuid = 1;</code>
     */
    private $uuid = '';
    /**
     * <code>string path = 2;</code>
     */
    private $path = '';
    /**
     * <code>string hostname = 3;</code>
     */
    private $hostname = '';

    public function __construct() {
        \GPBMetadata\File::initOnce();
        parent::__construct();
    }

    /**
     * <code>string uuid = 1;</code>
     */
    public function getUuid()
    {
        return $this->uuid;
    }

    /**
     * <code>string uuid = 1;</code>
     */
    public function setUuid($var)
    {
        GPBUtil::checkString($var, True);
        $this->uuid = $var;
    }

    /**
     * <code>string path = 2;</code>
     */
    public function getPath()
    {
        return $this->path;
    }

    /**
     * <code>string path = 2;</code>
     */
    public function setPath($var)
    {
        GPBUtil::checkString($var, True);
        $this->path = $var;
    }

    /**
     * <code>string hostname = 3;</code>
     */
    public function getHostname()
    {
        return $this->hostname;
    }

    /**
     * <code>string hostname = 3;</code>
     */
    public function setHostname($var)
    {
        GPBUtil::checkString($var, True);
        $this->hostname = $var;
    }

}

