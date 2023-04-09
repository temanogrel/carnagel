<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: common.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>pb.FileData</code>
 */
class FileData extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string uuid = 1;</code>
     */
    private $uuid = '';
    /**
     * Generated from protobuf field <code>uint64 externalId = 2;</code>
     */
    private $externalId = 0;
    /**
     * Generated from protobuf field <code>.pb.FileType type = 3;</code>
     */
    private $type = 0;
    /**
     * Generated from protobuf field <code>string hostname = 4;</code>
     */
    private $hostname = '';
    /**
     * Generated from protobuf field <code>string path = 5;</code>
     */
    private $path = '';
    /**
     * Generated from protobuf field <code>string upstoreHash = 6;</code>
     */
    private $upstoreHash = '';
    /**
     * Generated from protobuf field <code>string originalFilename = 7;</code>
     */
    private $originalFilename = '';
    /**
     * Generated from protobuf field <code>string checksum = 14;</code>
     */
    private $checksum = '';
    /**
     * Generated from protobuf field <code>bool pendingUpload = 8;</code>
     */
    private $pendingUpload = false;
    /**
     * Generated from protobuf field <code>bool pendingDeletion = 9;</code>
     */
    private $pendingDeletion = false;
    /**
     * Generated from protobuf field <code>uint64 size = 10;</code>
     */
    private $size = 0;
    /**
     * Generated from protobuf field <code>.pb.Struct meta = 11;</code>
     */
    private $meta = null;
    /**
     * Generated from protobuf field <code>string createdAt = 12;</code>
     */
    private $createdAt = '';
    /**
     * Generated from protobuf field <code>string updatedAt = 13;</code>
     */
    private $updatedAt = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $uuid
     *     @type int|string $externalId
     *     @type int $type
     *     @type string $hostname
     *     @type string $path
     *     @type string $upstoreHash
     *     @type string $originalFilename
     *     @type string $checksum
     *     @type bool $pendingUpload
     *     @type bool $pendingDeletion
     *     @type int|string $size
     *     @type \Pb\Struct $meta
     *     @type string $createdAt
     *     @type string $updatedAt
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Common::initOnce();
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
     * Generated from protobuf field <code>uint64 externalId = 2;</code>
     * @return int|string
     */
    public function getExternalId()
    {
        return $this->externalId;
    }

    /**
     * Generated from protobuf field <code>uint64 externalId = 2;</code>
     * @param int|string $var
     * @return $this
     */
    public function setExternalId($var)
    {
        GPBUtil::checkUint64($var);
        $this->externalId = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>.pb.FileType type = 3;</code>
     * @return int
     */
    public function getType()
    {
        return $this->type;
    }

    /**
     * Generated from protobuf field <code>.pb.FileType type = 3;</code>
     * @param int $var
     * @return $this
     */
    public function setType($var)
    {
        GPBUtil::checkEnum($var, \Pb\FileType::class);
        $this->type = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string hostname = 4;</code>
     * @return string
     */
    public function getHostname()
    {
        return $this->hostname;
    }

    /**
     * Generated from protobuf field <code>string hostname = 4;</code>
     * @param string $var
     * @return $this
     */
    public function setHostname($var)
    {
        GPBUtil::checkString($var, True);
        $this->hostname = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string path = 5;</code>
     * @return string
     */
    public function getPath()
    {
        return $this->path;
    }

    /**
     * Generated from protobuf field <code>string path = 5;</code>
     * @param string $var
     * @return $this
     */
    public function setPath($var)
    {
        GPBUtil::checkString($var, True);
        $this->path = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string upstoreHash = 6;</code>
     * @return string
     */
    public function getUpstoreHash()
    {
        return $this->upstoreHash;
    }

    /**
     * Generated from protobuf field <code>string upstoreHash = 6;</code>
     * @param string $var
     * @return $this
     */
    public function setUpstoreHash($var)
    {
        GPBUtil::checkString($var, True);
        $this->upstoreHash = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string originalFilename = 7;</code>
     * @return string
     */
    public function getOriginalFilename()
    {
        return $this->originalFilename;
    }

    /**
     * Generated from protobuf field <code>string originalFilename = 7;</code>
     * @param string $var
     * @return $this
     */
    public function setOriginalFilename($var)
    {
        GPBUtil::checkString($var, True);
        $this->originalFilename = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string checksum = 14;</code>
     * @return string
     */
    public function getChecksum()
    {
        return $this->checksum;
    }

    /**
     * Generated from protobuf field <code>string checksum = 14;</code>
     * @param string $var
     * @return $this
     */
    public function setChecksum($var)
    {
        GPBUtil::checkString($var, True);
        $this->checksum = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>bool pendingUpload = 8;</code>
     * @return bool
     */
    public function getPendingUpload()
    {
        return $this->pendingUpload;
    }

    /**
     * Generated from protobuf field <code>bool pendingUpload = 8;</code>
     * @param bool $var
     * @return $this
     */
    public function setPendingUpload($var)
    {
        GPBUtil::checkBool($var);
        $this->pendingUpload = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>bool pendingDeletion = 9;</code>
     * @return bool
     */
    public function getPendingDeletion()
    {
        return $this->pendingDeletion;
    }

    /**
     * Generated from protobuf field <code>bool pendingDeletion = 9;</code>
     * @param bool $var
     * @return $this
     */
    public function setPendingDeletion($var)
    {
        GPBUtil::checkBool($var);
        $this->pendingDeletion = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>uint64 size = 10;</code>
     * @return int|string
     */
    public function getSize()
    {
        return $this->size;
    }

    /**
     * Generated from protobuf field <code>uint64 size = 10;</code>
     * @param int|string $var
     * @return $this
     */
    public function setSize($var)
    {
        GPBUtil::checkUint64($var);
        $this->size = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>.pb.Struct meta = 11;</code>
     * @return \Pb\Struct
     */
    public function getMeta()
    {
        return $this->meta;
    }

    /**
     * Generated from protobuf field <code>.pb.Struct meta = 11;</code>
     * @param \Pb\Struct $var
     * @return $this
     */
    public function setMeta($var)
    {
        GPBUtil::checkMessage($var, \Pb\Struct::class);
        $this->meta = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string createdAt = 12;</code>
     * @return string
     */
    public function getCreatedAt()
    {
        return $this->createdAt;
    }

    /**
     * Generated from protobuf field <code>string createdAt = 12;</code>
     * @param string $var
     * @return $this
     */
    public function setCreatedAt($var)
    {
        GPBUtil::checkString($var, True);
        $this->createdAt = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string updatedAt = 13;</code>
     * @return string
     */
    public function getUpdatedAt()
    {
        return $this->updatedAt;
    }

    /**
     * Generated from protobuf field <code>string updatedAt = 13;</code>
     * @param string $var
     * @return $this
     */
    public function setUpdatedAt($var)
    {
        GPBUtil::checkString($var, True);
        $this->updatedAt = $var;

        return $this;
    }

}

