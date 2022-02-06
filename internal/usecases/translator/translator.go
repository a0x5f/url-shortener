package translator

import (
	"context"
	"errors"
	netUrl "net/url"
	"strings"
	"url-shortener/internal/domain/interfaces"
	"url-shortener/internal/domain/model"
	pkgErrors "url-shortener/pkg/errors"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	base     = uint64(len(alphabet))
	hashLen  = 10

	ShortLinkDomain = "https://lnk.dev/"
)

var (
	ErrorInvalidHashLen       = errors.New("invalid hash len")
	ErrorInvalidHashCharacter = errors.New("invalid hash character")
)

type UrlTranslator struct {
	dbService interfaces.DataStorageService
}

func New(dbService interfaces.DataStorageService) *UrlTranslator {
	return &UrlTranslator{
		dbService: dbService,
	}
}

func (h *UrlTranslator) ShortenUrl(ctx context.Context, url string) (string, error) {
	if _, err := netUrl.ParseRequestURI(url); err != nil {
		return "", err
	}

	err := h.dbService.Insert(ctx, &model.Link{Url: url})
	if err != nil {
		return "", pkgErrors.Wrap(err, "Couldn't insert record to database")
	}

	link, err := h.dbService.Select(ctx, &model.Link{Url: url})
	if err != nil {
		return "", pkgErrors.Wrap(err, "Couldn't select record from database")
	}

	shortened := h.getShortLinkById(link.Id)
	return shortened, nil
}

func (h *UrlTranslator) ExtendUrl(ctx context.Context, url string) (string, error) {
	hash := url[len(ShortLinkDomain):]

	if len(hash) != hashLen {
		return "", ErrorInvalidHashLen
	}

	for _, c := range hash {
		if !strings.ContainsRune(alphabet, c) {
			return "", ErrorInvalidHashCharacter
		}
	}

	id := h.getLinkIdByHash(hash)

	link, err := h.dbService.Select(ctx, &model.Link{Id: id})
	if err != nil {
		return "", pkgErrors.Wrap(err, "Couldn't select record from database")
	}

	return link.Url, nil
}

func (h *UrlTranslator) getShortLinkById(id uint64) string {
	hash := ""

	for id != 0 {
		digit := id % base
		hash += string(alphabet[digit])
		id = id / base
	}

	hash = reverseString(hash)

	for len(hash) < hashLen {
		hash = string(alphabet[0]) + hash
	}

	return ShortLinkDomain + hash
}

func reverseString(s string) string {
	result := ""

	for _, c := range s {
		result = string(c) + result
	}

	return result
}

func (h *UrlTranslator) getLinkIdByHash(hash string) uint64 {
	id := uint64(0)
	pow := uint64(1)

	for i := 0; i < len(hash); i++ {
		char := string(hash[len(hash)-i-1])
		id += uint64(strings.Index(alphabet, char)) * pow
		pow *= base
	}

	return id
}
