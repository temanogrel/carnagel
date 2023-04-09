<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: file.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>pb.RecommendStorageResponse</code>
 */
class RecommendStorageResponse extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>.pb.StatusCode status = 1;</code>
     */
    private $status = 0;
    /**
     * Generated from protobuf field <code>string hostname = 2;</code>
     */
    private $hostname = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type int $status
     *     @type string $hostname
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\File::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>.pb.StatusCode status = 1;</code>
     * @return int
     */
    public function getStatus()
    {
        return $this->status;
    }

    /**
     * Generated from protobuf field <code>.pb.StatusCode status = 1;</code>
     * @param int $var
     * @return $this
     */
    public function setStatus($var)
    {
        GPBUtil::checkEnum($var, \Pb\StatusCode::class);
        $this->status = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string hostname = 2;</code>
     * @return string
     */
    public function getHostname()
    {
        return $this->hostname;
    }

    /**
     * Generated from protobuf field <code>string hostname = 2;</code>
     * @param string $var
     * @return $this
     */
    public function setHostname($var)
    {
        GPBUtil::checkString($var, True);
        $this->hostname = $var;

        return $this;
    }

}
