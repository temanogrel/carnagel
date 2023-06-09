<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: minion.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>pb.RelocateRequest</code>
 */
class RelocateRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string uuid = 1;</code>
     */
    private $uuid = '';
    /**
     * Generated from protobuf field <code>string targetHost = 2;</code>
     */
    private $targetHost = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $uuid
     *     @type string $targetHost
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Minion::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string uuid = 1;</code>
     * @return string
     */
    public function getUuid()
    {
        return $this->uuid;
    }

    /**
     * Generated from protobuf field <code>string uuid = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setUuid($var)
    {
        GPBUtil::checkString($var, True);
        $this->uuid = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string targetHost = 2;</code>
     * @return string
     */
    public function getTargetHost()
    {
        return $this->targetHost;
    }

    /**
     * Generated from protobuf field <code>string targetHost = 2;</code>
     * @param string $var
     * @return $this
     */
    public function setTargetHost($var)
    {
        GPBUtil::checkString($var, True);
        $this->targetHost = $var;

        return $this;
    }

}

