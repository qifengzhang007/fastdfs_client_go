//  fastdfs  C语言源代码有关资源鉴权函数的实现

int fdfs_http_gen_token(const BufferInfo *secret_key, const char *file_id,
		const time_t timestamp, char *token)
{
	char buff[256 + 64];
	unsigned char digit[16];
	int id_len;
	int total_len;

	id_len = strlen(file_id);
	if (id_len + secret_key->length + 12 > sizeof(buff))
	{
		return ENOSPC;
	}

	memcpy(buff, file_id, id_len);
	total_len = id_len;
	memcpy(buff + total_len, secret_key->buff, secret_key->length);
	total_len += secret_key->length;
	total_len += fc_itoa(timestamp, buff + total_len);

	my_md5_buffer(buff, total_len, digit);
	bin2hex((char *)digit, 16, token);
	return 0;
}

int fdfs_http_check_token(const BufferInfo *secret_key, const char *file_id, \
		const time_t timestamp, const char *token, const int ttl)
{
	char true_token[33];
	int result;
	int token_len;

	token_len = strlen(token);
	if (token_len != 32)
	{
		return EINVAL;
	}

	if ((timestamp != 0) && (time(NULL) - timestamp > ttl))
	{
		return ETIMEDOUT;
	}

	if ((result=fdfs_http_gen_token(secret_key, file_id, \
			timestamp, true_token)) != 0)
	{
		return result;
	}

	return (memcmp(token, true_token, 32) == 0) ? 0 : EPERM;
}
